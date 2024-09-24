use std::time;

pub fn get_heartbeat() -> String {
    let current_timestamp = time::SystemTime::now()
        .duration_since(time::UNIX_EPOCH)
        .expect("SystemTimeError comparing current time and unix timestamp")
        .as_millis();

    let mut message = "PING ".to_string();
    message.push_str(&current_timestamp.to_string());
    message.push_str("\n");

    message
}