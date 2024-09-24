pub enum Command {
    UNKNOWN,
    CONNECT,
    CLOSE
}

impl Command {
    pub fn parse(message: &str) -> Command {
        println!("[command] message received: {}", message);
        
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