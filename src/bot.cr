require "marmot"

require "./aoc"
require "./matrix"
require "./repository"

class Rudolphe::Bot
  @config : Config
  @db_leaderboard : Leaderboard

  def initialize
    @repository = Repository.new
    @config = @repository.get_config
    @matrix = Matrix.new(@config, @repository)
    @matrix.set_sync_task
    @aoc = Aoc.new(@config)

    @db_leaderboard = @repository.get_leaderboard

    @quiet = false
    @leaderboard_task = Marmot.repeat(15.minutes, true) { check_leaderboard }

    # The new puzzle is out on midnight UTC-5, which is 6:00 at UTC+1
    @puzzle_task = Marmot.cron(6, 0) { send_puzzle_link }

    Marmot.run
  end

  def check_leaderboard : Nil
    Log.info { "Checking leaderboard" }
    @aoc.get_leaderboard.try &.users.each do |user_id, user|
      if db_user = @db_leaderboard.users[user_id]?
        new_points = nil
        if user.local_score != db_user.local_score
          new_points = user.local_score - db_user.local_score
          @repository.save_user_local_score(user)
        end

        user.days.each do |day, parts|
          if db_day = db_user.days[day]?
            parts.each_key do |part|
              if !db_day.has_key?(part)
                @matrix.send("#{user.name} vient juste de compléter la partie #{part} du jour #{day} (+#{new_points} points)")
                @repository.save_user_part(user_id, day, part, parts[part])
              end
            end
          else
            if parts.size == 2
              parts_msg = "les 2 parties"
            else
              parts_msg = "la partie #{parts.first_key}"
            end

            @matrix.send("#{user.name} vient juste de compléter #{parts_msg} du jour #{day} (+#{new_points} points)")
            @repository.save_user_day(user_id, day, parts)
          end
        end
      else
        @matrix.send("Un nouveau concurrent entre dans la place, bienvenue à #{user.name} !")
        @repository.save_user(user)
      end

      @db_leaderboard.users[user_id] = user
    end

    reschedule_tasks
  end

  def reschedule_tasks
    now = Time.utc
    if !@quiet && (now.month != 12 || !(1..25).includes?(now.day))
      Log.info { "Going into quiet mode" }
      @quiet = true
      @puzzle_task.cancel
      @leaderboard_task.cancel
      @leaderboard_task = Marmot.repeat(1.hour) { check_leaderboard }
    elsif @quiet && now.month == 12 && (1..25).includes?(now.day)
      Log.info { "Going into competiton mode" }
      @quiet = false
      @puzzle_task = Marmot.cron(6, 0) { send_puzzle_link }
      @leaderboard_task.cancel
      @leaderboard_task = Marmot.repeat(15.minutes) { check_leaderboard }
    end
  end

  def send_puzzle_link
    @matrix.send("Nouveau puzzle : https://adventofcode.com/#{Time.local.year}/day/#{Time.local.day}")
    @matrix.send_leaderboard
  end
end
