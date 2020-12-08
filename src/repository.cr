require "db"
require "sqlite3"

require "./config"
require "./leaderboard"

class Rudolphe::Repository
  Log = Rudolphe::Log.for(self)

  @db : DB::Database

  def initialize
    @db = DB.open("sqlite3:./data.db")
  end

  def get_config : Config
    @db.query("select access_token, room, session_cookie, leaderboard_url from config") do |rs|
      rs.each do
        access_token = rs.read(String)
        room = rs.read(String)
        session_cookie = rs.read(String)
        leaderboard_url = rs.read(String)
        return Config.new(access_token, room, session_cookie, leaderboard_url)
      end
    end

    Log.fatal { "missing configuration" }
    exit(1)
  end

  def get_leaderboard : Leaderboard
    users = Hash(String, User).new
    @db.query("
      select u.id, u.name, d.day, d.part, d.get_star_ts
      from users u
      left outer join days d on
        u.id = d.user_id
    ") do |rs|
      rs.each do
        user_id = rs.read(String)
        user_name = rs.read(String)
        day = rs.read(Int64?).try &.to_u8
        part = rs.read(Int64?).try &.to_u8
        get_star_ts = rs.read(String?)

        if user = users[user_id]?
          if !day.nil? && !part.nil? && !get_star_ts.nil?
            if user_day = user.days[day]?
              user_day[part] = Star.new(get_star_ts)
            else
              user.days[day] = {part => Star.new(get_star_ts)}
            end
          end
        else
          if !day.nil? && !part.nil? && !get_star_ts.nil?
            days = {day => {part => Star.new(get_star_ts)}}
          else
            days = Hash(UInt8, Hash(UInt8, Star)).new
          end

          users[user_id] = User.new(
            days: days,
            id: user_id,
            name: user_name,
          )
        end
      end
    end

    Leaderboard.new(users)
  end

  def save_user(user) : Nil
    @db.exec(
      "insert into users (id, name) values (?, ?)",
      user.id,
      user.name,
    )

    user.days.each do |day, parts|
      save_user_day(user.id, day, parts)
    end
  end

  def save_user_day(user_id, day, parts) : Nil
    parts.each do |part, star|
      save_user_part(user_id, day, part, star)
    end
  end

  def save_user_part(user_id, day, part, star)
    @db.exec(
      "insert into days (user_id, day, part, get_star_ts) values (?, ?, ?, ?)",
      user_id,
      day.to_i,
      part.to_i,
      star.get_star_ts,
    )
  end
end
