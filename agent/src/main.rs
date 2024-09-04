use std::{thread, time};

fn main() {
    let heartbeat_interval = time::Duration::from_millis(1000);

    loop {
        let now = time::SystemTime::now()
            .duration_since(time::UNIX_EPOCH)
            .expect("SystemTimeError comparing current time and unix timestamp")
            .as_millis();

        println!("PING {now}");
        thread::sleep(heartbeat_interval);
    }
}
