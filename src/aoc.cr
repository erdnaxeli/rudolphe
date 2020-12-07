require "http/client"
require "http/cookie"
require "http/headers"

class Rudolphe::Aoc
  Log = Rudolphe::Log.for(self)

  def initialize(@config : Config)
  end

  def get_leaderboard : Leaderboard?
    headers = HTTP::Headers.new
    cookies = HTTP::Cookies.new
    cookies << HTTP::Cookie.new("session", @config.session_cookie)
    cookies.add_request_headers(headers)

    response = HTTP::Client.get(@config.leaderboard_url, headers)
    if !response.success?
      Log.warn &.emit(
        "Error while getting leaderboard",
        status_code: response.status_code,
        body: response.body
      )
      return
    end

    Leaderboard.from_json(response.body)
  end
end
