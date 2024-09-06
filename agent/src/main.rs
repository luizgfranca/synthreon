use std::{io::{BufRead, BufReader, Write}, net::{TcpListener, TcpStream}, thread, time};

static HEARTBEAT_INTERVAL: std::time::Duration = time::Duration::from_millis(1000);

enum Command {
    UNKNOWN,
    CONNECT,
    CLOSE
}

impl Command {
    pub fn from(message: &str) -> Command {
        if message.is_empty() {
            Command::CLOSE
        } else {
            match message {
                "connect" => Command::CONNECT,
                _ => Command::UNKNOWN
            }
        }
    }
}

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

fn transmit_heartbeat(mut stream: TcpStream) {
    loop {
        let current_timestamp = time::SystemTime::now()
            .duration_since(time::UNIX_EPOCH)
            .expect("SystemTimeError comparing current time and unix timestamp")
            .as_millis();

        let mut message = "PING ".to_string();
        message.push_str(&current_timestamp.to_string());
        message.push_str("\n");
        print!("{message}");

        let result = stream.write_all(message.as_bytes());
        if result.is_err() {
            println!("error sending heartbeat\n");
            break;
        }

        thread::sleep(HEARTBEAT_INTERVAL);
    }
}

fn handle_connection(stream: TcpStream) {
    println!("new connection");
    let command = Command::from(&wait_for_message(&stream));
    match command {
        Command::CONNECT => transmit_heartbeat(stream),
        Command::UNKNOWN => println!("unknown command received"),
        Command::CLOSE => println!("close connection"),
    }
}

fn main() {
    let address = "127.0.0.1:25250";
    let listener = TcpListener::bind(address)
        .expect("failed to attach to {address}");

    for stream in listener.incoming() {
        let stream = stream.unwrap();
        handle_connection(stream);
    }

}
