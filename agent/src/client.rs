use std::{
    io::{ErrorKind, Read, Write}, net::TcpStream
};

use crate::command::Command;

pub enum ClientStatus {
    Connected,
    Subscribed,
    Disconnected,
}

pub struct Client {
    pub stream: TcpStream,
    pub status: ClientStatus,
}

impl Client {
    pub fn new(stream: TcpStream) -> Self {
        stream.set_nonblocking(true).expect("unable to set stream as nonBlocking");
        let client = Client {
            stream,
            status: ClientStatus::Connected,
        };

        client
    }

    pub fn process_new_commands(&mut self) {
        for command in self.get_new_commands() {
            match command {
                Command::CONNECT => self.set_subscribed(),
                Command::CLOSE => self.disconnect(),
                Command::UNKNOWN => println!("unknwon command received"),
            }
        }
    }

    pub fn send_message(&mut self, message: &str) {
        println!("[client] sending message: {}", message.trim());
        self.stream
            .write_all(message.as_bytes())
            .unwrap_or_else(move |_| {
                self.disconnect();
            });
    }

    // TODO: this approach is terible, and should be redone
    fn get_new_commands(& mut self) -> Vec<Command> {
        let mut commands = Vec::new();
        let mut buffer = [0];
        let mut message = String::new();
    
        loop {
            match self.stream.read(&mut buffer) {
                Err(error) => {
                    match error.kind() {
                        ErrorKind::WouldBlock => break,
                        _ => {
                            println!("error reading from stream: {}", error);
                            break;
                        }
                    }
                    
                }
                Ok(got) => {
                    match got {
                        0 => break,
                        1 => {
                            message.push_str(
                                String::from_utf8(
                                    Vec::from(&buffer)
                                ).expect("unable to convert string message")
                                .as_str()
                            );
                            
                            let last_char = message.chars()
                                .last()
                                .expect("string empty after being filled with character");

                            if last_char == '\n' || last_char == '\0' {
                                let command = Command::parse(&message[0..message.len() - 1]);
                                commands.push(command);
                                
                                // TODO: small memory leak here until the end of the function,
                                //       but currently it is the least of its problems
                                message = String::new();
                            }
                        },
                        _ => panic!("unexpected number of characters loaded")
                    }
                }
            }


        }

        commands
    }

    fn set_subscribed(&mut self) {
        self.status = ClientStatus::Subscribed;
        println!("[client] client subscribed");
    }

    fn disconnect(&mut self) {
        self.status = ClientStatus::Disconnected;
    }
}
