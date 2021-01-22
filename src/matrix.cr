require "caridina"
require "caridina/syncer"

require "marmot"

class Rudolphe::Matrix
  @first_sync = true

  def initialize(@config : Config, @repository : Rudolphe::Repository)
    @conn = Caridina::Connection.new(
      "https://cloud.cervoi.se",
      config.access_token,
    )
    @conn.join(@config.room)
  end

  def send(message, formatted_message : String? = nil) : Nil
    @conn.send_message(@config.room, message, formatted_message)
  end

  def send_leaderboard
    leaderboard = @repository.get_leaderboard
    msg = leaderboard.to_s
    send(msg, "<pre>#{msg}</pre>")
  end

  def set_sync_task : Nil
    channel = Channel(Caridina::Responses::Sync).new
    @conn.sync(channel)

    syncer = Caridina::Syncer.new
    syncer.on(Caridina::Events::Message) do |event|
      on_message(event)
    end

    Marmot.on(channel) do |task|
      response = task.as(Marmot::OnChannelTask).value
      if !response.nil?
        syncer.process_response(response)
        @first_sync = false if @first_sync
      end
    end
  end

  def on_message(event)
    # Skip the first sync messages as it can contains messages already read.
    return if @first_sync

    event = event.as(Caridina::Events::Message)
    room_id = event.room_id.not_nil!
    @conn.send_receipt(room_id, event.event_id)

    if event.sender != @conn.user_id && (message = event.content.as?(Caridina::Events::Message::Text))
      if message.body == "!lb"
        send_leaderboard
      end
    end
  end
end
