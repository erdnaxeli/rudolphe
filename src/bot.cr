require "marmot"

require "./aoc"
require "./matrix"
require "./repository"

class Rudolphe::Bot
  @config : Config
  @current_leaderboard : Leaderboard

  def initialize
    @repository = Repository.new
    @config = @repository.get_config
    @matrix = Matrix.new(@config, @repository)
    @matrix.set_sync_task
    @aoc = Aoc.new(@config)

    @current_leaderboard = @repository.get_leaderboard

    @quiet = false
    @leaderboard_task = Marmot.every(15.minutes, true) { check_leaderboard }

    # The new puzzle is out on midnight UTC-5, which is 6:00 at UTC+1
    @puzzle_task = Marmot.every(:day, hour: 6) { send_puzzle_link }

    Marmot.run
  end

  def check_leaderboard : Nil
    Log.info { "Checking leaderboard" }
    if leaderboard = @aoc.get_leaderboard
      @current_leaderboard.diff_to(leaderboard).try &.users.each do |user_id, user|
        if @current_leaderboard.users.has_key?(user_id)
          days = user.days.map do |day, parts|
            if parts.size == 2
              "le jour #{day}"
            else
              "la partie #{parts.first_key} du jour #{day}"
            end
          end

          msg = String.build do |str|
            str << user.name << " vient juste de compléter "
            if days.size == 1
              str << days[0]
            else
              days[...-2].each { |d| str << d << ", " }
              str << days[-2] << " et " << days[-1]
            end

            str << " (+" << user.local_score << " points)"
          end

          @matrix.send(msg)
        else
          msg = String.build do |str|
            str << "Un nouveau concurrent entre dans la place"
            if user.local_score > 1
              str << " avec " << user.local_score << " points"
            end

            str << ", bienvenue à " << user.name << " !"
          end
          @matrix.send(msg)
        end
      end

      @current_leaderboard = leaderboard
      @repository.save_leaderboard(leaderboard)
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
      @leaderboard_task = Marmot.every(1.hour) { check_leaderboard }
    elsif @quiet && now.month == 12 && (1..25).includes?(now.day)
      Log.info { "Going into competiton mode" }
      @quiet = false
      @puzzle_task = Marmot.every(:day, hour: 6) { send_puzzle_link }
      @leaderboard_task.cancel
      @leaderboard_task = Marmot.every(15.minutes) { check_leaderboard }
    end
  end

  def send_puzzle_link
    @matrix.send("Nouveau puzzle : https://adventofcode.com/#{Time.local.year}/day/#{Time.local.day}")
    @matrix.send_leaderboard
  end
end
