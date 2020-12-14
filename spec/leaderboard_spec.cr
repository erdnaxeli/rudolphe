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
end
