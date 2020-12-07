require "json"

struct Leaderboard
  include JSON::Serializable

  @[JSON::Field(key: "members")]
  getter users : Hash(String, User)

  def initialize(@users)
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
end

struct Star
  include JSON::Serializable

  getter get_star_ts : String

  def initialize(@get_star_ts)
  end
end
