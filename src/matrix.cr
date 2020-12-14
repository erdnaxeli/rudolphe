require "caridina"
require "caridina/syncer"

require "marmot"

class Rudolphe::Matrix
  @first_sync = true

  def initialize(@config : Config, @repository : Rudolphe::Repository)
    @conn = Caridina::ConnectionImpl.new(
      "https://cloud.cervoi.se",
      config.access_token,
    )
    @conn.join(@config.room)
  end

  def send(message, formatted_message : String? = nil) : Nil
    @conn.send_message(@config.room, message, formatted_message)
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
        msg = String.build do |str|
          leaderboard = @repository.get_leaderboard
          score_size = 0
          position_size = leaderboard.users.size.to_s.size

          users = leaderboard.users.values.sort do |user|
            score_size = {score_size, user.local_score.to_s.size}.max
            user.local_score.to_i
          end

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

        send(msg, "<pre>#{msg}</pre>")
      end
    end
  end
end
