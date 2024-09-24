use std::net::TcpListener;
use crate::client::Client;

pub struct ClientManager {
    listener: TcpListener,
    pub clients: Vec<Client>
}

impl ClientManager {
    pub fn new(listener: TcpListener) -> Self {
        listener.set_nonblocking(true)
            .expect("unable to set connection as nonblocking");
        
        Self {
            listener,
            clients: Vec::new(),
        }
    }

    pub fn check_new_connections(&mut self) {
        println!("[clientManager] polling for new connections");
        match self.listener.accept() {
            Ok((stream, _)) => self.clients.push(Client::new(stream)),
            Err(_) => ()
        }
    }

    pub fn process_new_commands(&mut self) {
        for client in self.clients.iter_mut() {
            client.process_new_commands();
        }
    }

    pub fn broadcast_to_subscribers(&mut self, message: &str) {
        for client in self.clients.iter_mut() {
            match client.status {
                crate::client::ClientStatus::Subscribed => client.send_message(message),
                _ => ()
            }
        }
    }
}
