require "json"

struct Leaderboard
  include JSON::Serializable

  @[JSON::Field(key: "members")]
  getter users : Hash(String, User)

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
end

struct User
  include JSON::Serializable

  @[JSON::Field(key: "completion_day_level")]
  getter days : Hash(UInt8, Hash(UInt8, Star)) = Hash(UInt8, Hash(UInt8, Star)).new
  getter id : String
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
end

struct Star
  include JSON::Serializable

  getter get_star_ts : String

  def initialize(@get_star_ts)
  end
end
