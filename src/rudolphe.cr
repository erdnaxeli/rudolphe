require "log"

require "./aoc"
require "./matrix"
require "./repository"

module Rudolphe
  VERSION = "0.1.0"

  Log = ::Log.for(self)

  def self.fly
    repository = Repository.new
    config = repository.get_config
    matrix = Matrix.new(config)
    aoc = Aoc.new(config)

    db_leaderboard = repository.get_leaderboard

    loop do
      aoc.get_leaderboard.try &.users.each do |user_id, user|
        if db_user = db_leaderboard.users[user_id]?
          user.days.each do |day, parts|
            if db_day = db_user.days[day]?
              parts.each_key do |part|
                if !db_day.has_key?(part)
                  matrix.send("#{user.name} vient juste de compléter la partie #{part} du jour #{day}")
                  db_day[part] = parts[part]
                  repository.save_user_part(user_id, day, part, parts[part])
                end
              end
            else
              if parts.size == 2
                parts_msg = "les 2 parties"
              else
                parts_msg = "la partie #{parts.first_key}"
              end

              matrix.send("#{user.name} vient juste de compléter #{parts_msg} pour le jour #{day}")
              repository.save_user_day(user_id, day, parts)
              db_user.days[day] = parts
            end
          end
        else
          matrix.send("Un nouveau concurrent entre dans la place, bienvenue à #{user.name} !")
          repository.save_user(user)
          db_leaderboard.users[user_id] = user
        end
      end

      Log.info { "Sleeping for 15 minutes" }
      sleep 15.minutes
    end
  end
end
