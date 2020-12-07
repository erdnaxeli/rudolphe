require "caridina"

class Rudolphe::Matrix
  def initialize(@config : Config)
    @conn = Caridina::ConnectionImpl.new(
      "https://cloud.cervoi.se",
      config.access_token,
    )
    @conn.join(@config.room)
  end

  def send(message) : Nil
    @conn.send_message(@config.room, message)
  end
end
