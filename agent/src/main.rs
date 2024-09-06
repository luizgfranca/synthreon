use std::{io::{BufRead, BufReader, Write}, net::{TcpListener, TcpStream}, thread, time};

fn wait_for_message(stream: &TcpStream) -> String {
    let reader = BufReader::new(stream);
            
    let messages: Vec<_> = reader
        .lines()
        .map(|message| message.unwrap())
        .take(1)
        .collect();

    let message = match messages.get(0) {
        Some(m) => m,
        None => ""
    };

    message.to_string()
}

fn handle_connection(mut stream: TcpStream) {
    println!("new connection");
    
    loop {
        let message = wait_for_message(&stream);
        if message.is_empty() {
            break;
        }

        println!("{message}");
        stream.write_all(message.as_bytes()).expect("error responding");
    }
}

fn main() {
    let heartbeat_interval = time::Duration::from_millis(1000);

    // loop {
        let now = time::SystemTime::now()
            .duration_since(time::UNIX_EPOCH)
            .expect("SystemTimeError comparing current time and unix timestamp")
            .as_millis();

        println!("PING {now}");
        thread::sleep(heartbeat_interval);
    // }

    let address = "127.0.0.1:25250";
    let listener = TcpListener::bind(address)
        .expect("failed to attach to {address}");

    for stream in listener.incoming() {
        let stream = stream.unwrap();
        handle_connection(stream);
    }

}
