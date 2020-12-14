require "log"

require "marmot"

require "./aoc"
require "./matrix"
require "./repository"

module Rudolphe
  VERSION = "0.1.0"

  Log = ::Log.for(self)

  def self.fly
    repository = Repository.new
    config = repository.get_config
    matrix = Matrix.new(config, repository)
    matrix.set_sync_task
    aoc = Aoc.new(config)

    db_leaderboard = repository.get_leaderboard

    Marmot.repeat(15.minutes, true) do
      Log.info { "Checking leaderboard" }
      aoc.get_leaderboard.try &.users.each do |user_id, user|
        if db_user = db_leaderboard.users[user_id]?
          if user.local_score != db_user.local_score
            repository.save_user_local_score(user)
          end

          user.days.each do |day, parts|
            if db_day = db_user.days[day]?
              parts.each_key do |part|
                if !db_day.has_key?(part)
                  matrix.send("#{user.name} vient juste de compléter la partie #{part} du jour #{day}")
                  repository.save_user_part(user_id, day, part, parts[part])
                end
              end
            else
              if parts.size == 2
                parts_msg = "les 2 parties"
              else
                parts_msg = "la partie #{parts.first_key}"
              end

              matrix.send("#{user.name} vient juste de compléter #{parts_msg} du jour #{day}")
              repository.save_user_day(user_id, day, parts)
            end
          end
        else
          matrix.send("Un nouveau concurrent entre dans la place, bienvenue à #{user.name} !")
          repository.save_user(user)
        end

        db_leaderboard.users[user_id] = user
      end
    end

    # The new puzzle is out on midnight UTC-5, which is 6:00 at UTC+1
    Marmot.cron(6, 0) do
      matrix.send("Nouveau puzzle : https://adventofcode.com/2020/day/#{Time.local.day}")
    end

    Marmot.run
  end
end
