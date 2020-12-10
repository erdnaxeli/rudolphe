require "../src/scheduler"

describe Rudolphe::Scheduler do
  describe "#repeat" do
    it "schedules a new task that repeats" do
      tempfile = File.tempfile.path

      task = Rudolphe::Scheduler.repeat(50.milliseconds) { File.write(tempfile, "I am inside the task\n", mode: "a") }
      spawn Rudolphe::Scheduler.run

      sleep 55.milliseconds
      File.read_lines(tempfile).size.should eq(1)
      sleep 50.milliseconds
      File.read_lines(tempfile).size.should eq(2)
      sleep 50.milliseconds
      File.read_lines(tempfile).size.should eq(3)

      task.cancel
    end
  end

  describe "#cron" do
    it "schedules a new task" do
      tempfile = File.tempfile.path

      time = Time.local.at_beginning_of_minute + 1.minute
      tasks = Rudolphe::Scheduler.cron(time.hour, time.minute) { File.write(tempfile, "\n", mode: "a") }
      spawn Rudolphe::Scheduler.run

      sleep (time - Time.local + 5.milliseconds)
      File.read_lines(tempfile).size.should eq(1)
    end
  end
end
