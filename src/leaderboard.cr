require "json"

struct Leaderboard
  include JSON::Serializable

  @[JSON::Field(key: "members")]
  getter users : Hash(UInt64, User)

  def initialize(@users)
  end

  def to_s
    score_size = 0
    position_size = @users.size.to_s.size

    users = @users.values.sort_by do |user|
      score_size = {score_size, user.local_score.to_s.size}.max
      -user.local_score.to_i
    end

    String.build do |str|
      # header tens
      (1..position_size + score_size + 3).each { str << ' ' }
      (1..25).each do |d|
        if d < 10
          str << ' '
        elsif d < 20
          str << 1
        else
          str << 2
        end
      end
      str << '\n'

      # header units
      (1..position_size + score_size + 3).each { str << ' ' }
      (1..25).each do |d|
        str << d % 10
      end
      str << '\n'

      users.each_with_index(1) do |user, i|
        # position
        str << sprintf("%#{position_size}d", i) << ") "

        # score
        str << sprintf("%#{score_size}d", user.local_score) << ' '

        # stars
        (1..25).each do |d|
          day = user.days[d]?
          if day.nil?
            str << ' '
          elsif day.size == 1
            str << '+'
          else
            str << '×'
          end
        end

        # username
        str << ' ' << user.name_without_hl
        str << '\n' if i < users.size
      end
    end
  end

  def diff_to(other : Leaderboard) : Leaderboard?
    diff_users = Hash(UInt64, User).new
    other.users.each do |other_id, other_user|
      if user = @users[other_id]?
        if diff_user = user.diff_to(other_user)
          diff_users[other_id] = diff_user
        end
      else
        diff_users[other_id] = other_user
      end
    end

    if diff_users.size == 0
      nil
    else
      Leaderboard.new(users: diff_users)
    end
  end
end

struct User
  include JSON::Serializable

  @[JSON::Field(key: "completion_day_level")]
  getter days = Hash(UInt8, Hash(UInt8, Star)).new
  getter id : UInt64
  getter local_score : UInt16
  @name : String?

  def initialize(@days, @id, @name, @local_score = 0_u16)
  end

  def name : String
    @name || "anon #{@id}"
  end

  def name_without_hl : String
    # It joins with a "zero width non-joiner" char.
    name.split("").join('\u200C')
  end

  def diff_to(other : User) : User?
    new_points = other.local_score - @local_score
    diff_days = Hash(UInt8, Hash(UInt8, Star)).new { |h, k| h[k] = Hash(UInt8, Star).new }

    other.days.each do |other_day, other_parts|
      if parts = @days[other_day]?
        other_parts.each do |other_part, other_star|
          if !parts.has_key?(other_part)
            diff_days[other_day][other_part] = other_star
          end
        end
      else
        diff_days[other_day] = other_parts
      end
    end

    if diff_days.size == 0
      nil
    else
      User.new(days: diff_days, id: @id, name: @name, local_score: new_points)
    end
  end
end

struct Star
  include JSON::Serializable

  getter get_star_ts : Int32

  def initialize(@get_star_ts)
  end
end
