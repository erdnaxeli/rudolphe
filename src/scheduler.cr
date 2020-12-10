module Rudolphe::Scheduler
  @@tasks = Array(Task).new

  abstract class Task
    @canceled = false
    @callback : Proc(Nil) = ->{}

    getter tick = Channel(Task).new

    abstract def start : Nil

    def cancel
      @canceled = true
    end

    def run : Nil
      @callback.call
    end
  end

  class Repeat < Task
    def initialize(@span : Time::Span, @callback : Proc(Nil))
    end

    def start : Nil
      spawn do
        while !@canceled
          sleep @span
          @tick.send(self)
        end
      end
    end
  end

  class Cron < Task
    def initialize(@hour : Int32, @minute : Int32, @callback : Proc(Nil))
    end

    def start : Nil
      spawn do
        while !@canceled
          sleep span
          @tick.send(self)
        end
      end
    end

    private def span
      # We want the next minute, we skip the current one.
      time = Time.local.at_beginning_of_minute + 1.minute

      if time.minute < @minute
        time += (@minute - time.minute).minute
      elsif time.hour > @minute
        time += (60 - time.minute + @minute).minute
      end

      if time.hour < @hour
        time += (@hour - time.hour).hour
      elsif time.hour > @hour
        time += (24 - time.hour + @hour).hour
      end

      time - Time.local
    end
  end

  extend self

  # Runs a task every given *span*.
  def repeat(span : Time::Span, &block) : Task
    task = Repeat.new(span, block)
    @@tasks << task
    task
  end

  # Runs a task every day at *hour* and *minute*.
  def cron(hour, minute, &block) : Task
    task = Cron.new(hour, minute, block)
    @@tasks << task
    task
  end

  #  Starts scheduling the tasks.
  #
  # This blocks forever.
  def run
    @@tasks.map(&.start)

    loop do
      _, m = Channel.select(@@tasks.map(&.tick.receive_select_action))
      m.run
    end
  end
end
