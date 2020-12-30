require "./spec_helper"
require "../src/leaderboard"

Spectator.describe Leaderboard do
  let(json) do
    %(
      {
        "owner_id":"756077",
        "event":"2020",
        "members":{
          "756077":{
            "global_score":0,
            "completion_day_level":{
              "2":{
                "1":{"get_star_ts":"1607203578"},
                "2":{"get_star_ts":"1607205118"}
              },
              "6":{
                "1":{"get_star_ts":"1607284003"},
                "2":{"get_star_ts":"1607294072"}
              }
            },
            "local_score":42,
            "id":"756077",
            "last_star_ts":"1607338189",
            "name":null,
            "stars":14
          },
          "699989":{
            "local_score":104,
            "completion_day_level":{
              "5":{"1":{"get_star_ts":"1607193555"}}
            },
            "global_score":0,
            "id":"699989",
            "last_star_ts":"1607346014",
            "name":"EEEEEE",
            "stars":14
          }
        }
      }
    )
  end

  describe "#from_json" do
    it "deserializes an actual leaderboard" do
      leaderboard = Leaderboard.from_json(json)

      expect(leaderboard.users.size).to eq(2)
      expect(leaderboard.users["756077"].id).to eq("756077")
      expect(leaderboard.users["756077"].local_score).to eq(42)
      expect(leaderboard.users["756077"].name).to eq("anon 756077")

      expect(leaderboard.users["699989"].id).to eq("699989")
      expect(leaderboard.users["699989"].local_score).to eq(104)
      expect(leaderboard.users["699989"].name).to eq("EEEEEE")
    end
  end

  describe "#to_s" do
    it "returns a string representing a leaderboard" do
      leaderboard = Leaderboard.from_json(json)

      expect(leaderboard.to_s).to eq(
        %{                1111111111222222
       1234567890123456789012345
1) 104     +                     E\u{200C}E\u{200C}E\u{200C}E\u{200C}E\u{200C}E
2)  42  ×   ×                    a\u{200C}n\u{200C}o\u{200C}n\u{200C} \u{200C}7\u{200C}5\u{200C}6\u{200C}0\u{200C}7\u{200C}7}
      )
    end
  end

  describe "#diff_to" do
    it "returns nil when there is no diff" do
      leaderboard1 = Leaderboard.from_json(json)
      leaderboard2 = Leaderboard.from_json(json)

      diff = leaderboard1.diff_to(leaderboard2)
      expect(diff).to be_nil
    end

    it "returns a leaderboard with the new user when there is a new user" do
      leaderboard1 = Leaderboard.new(
        users: {
          "user1" => User.new(
            id: "user1",
            name: "User 1",
            days: Hash(UInt8, Hash(UInt8, Star)).new
          ),
        }
      )
      leaderboard2 = Leaderboard.new(
        users: {
          "user1" => User.new(
            id: "user1",
            name: "User 1",
            days: Hash(UInt8, Hash(UInt8, Star)).new
          ),
          "user2" => User.new(
            id: "user2",
            name: "User 2",
            days: {
              1_u8 => {1_u8 => Star.new("123"), 2_u8 => Star.new("123")},
            }
          ),
        }
      )

      diff = leaderboard1.diff_to(leaderboard2)
      expect(diff).to ne(nil)
      diff = diff.not_nil!

      expect(diff.users.size).to eq(1)
      expect(diff.users.first_key).to eq("user2")
    end

    it "returns a leaderboard with a user when the user has a diff" do
      leaderboard1 = Leaderboard.new(
        users: {
          "user1" => User.new(
            id: "user1",
            name: "User 1",
            days: Hash(UInt8, Hash(UInt8, Star)).new
          ),
        }
      )
      leaderboard2 = Leaderboard.new(
        users: {
          "user1" => User.new(
            id: "user1",
            name: "User 1",
            days: {1_u8 => {1_u8 => Star.new("123")}}
          ),
        }
      )

      diff = leaderboard1.diff_to(leaderboard2)
      expect(diff).to ne(nil)
      diff = diff.not_nil!

      expect(diff.users.size).to eq(1)
      expect(diff.users.first_key).to eq("user1")
    end
  end
end

Spectator.describe User do
  describe "#diff_to" do
    it "returns nil when there is not diff" do
      user1 = User.new(
        id: "User1",
        name: "User 1",
        days: Hash(UInt8, Hash(UInt8, Star)).new,
      )

      diff = user1.diff_to(user1)
      expect(diff).to be(nil)
    end

    it "returns a user with new days" do
      user1 = User.new(
        id: "User1",
        name: "User 1",
        days: Hash(UInt8, Hash(UInt8, Star)).new,
      )
      user2 = User.new(
        id: "User1",
        name: "User 1",
        days: {1_u8 => {1_u8 => Star.new("123")}},
      )

      diff = user1.diff_to(user2)
      expect(diff).to ne(nil)
      diff = diff.not_nil!

      expect(diff.id).to eq("User1")
      expect(diff.name).to eq("User 1")
      expect(diff.days.size).to eq(1)
      expect(diff.days.first_key).to eq(1)
      expect(diff.days[1].size).to eq(1)
      expect(diff.days[1].first_key).to eq(1)
      expect(diff.days[1][1].get_star_ts).to eq("123")
    end

    it "returns a user with new parts" do
      user1 = User.new(
        id: "User1",
        name: "User 1",
        days: {1_u8 => {1_u8 => Star.new("123")}},
      )
      user2 = User.new(
        id: "User1",
        name: "User 1",
        days: {
          1_u8 => {
            1_u8 => Star.new("123"),
            2_u8 => Star.new("234"),
          },
        },
      )

      diff = user1.diff_to(user2)
      expect(diff).to ne(nil)
      diff = diff.not_nil!

      expect(diff.id).to eq("User1")
      expect(diff.name).to eq("User 1")
      expect(diff.days.size).to eq(1)
      expect(diff.days.first_key).to eq(1)
      expect(diff.days[1].size).to eq(1)
      expect(diff.days[1].first_key).to eq(2)
      expect(diff.days[1][2].get_star_ts).to eq("234")
    end

    it "returns both new days and new parts" do
      user1 = User.new(
        id: "User1",
        name: "User 1",
        days: {1_u8 => {1_u8 => Star.new("123")}},
      )
      user2 = User.new(
        id: "User1",
        name: "User 1",
        days: {
          1_u8 => {
            1_u8 => Star.new("123"),
            2_u8 => Star.new("234"),
          },
          2_u8 => {1_u8 => Star.new("345")},
        },
      )

      diff = user1.diff_to(user2)
      expect(diff).to ne(nil)
      diff = diff.not_nil!

      expect(diff.id).to eq("User1")
      expect(diff.name).to eq("User 1")
      expect(diff.days.size).to eq(2)
      expect(diff.days.keys).to eq([1, 2])
      expect(diff.days[1].size).to eq(1)
      expect(diff.days[1].first_key).to eq(2)
      expect(diff.days[1][2].get_star_ts).to eq("234")
      expect(diff.days[2].size).to eq(1)
      expect(diff.days[2].first_key).to eq(1)
      expect(diff.days[2][1].get_star_ts).to eq("345")
    end
  end
end
