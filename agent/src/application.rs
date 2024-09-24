use std::time;
use std::net::TcpListener;

use crate::client_manager::ClientManager;
use crate::heartbeat;

static HEARTBEAT_INTERVAL: std::time::Duration = time::Duration::from_millis(1000);

pub fn serve_application(address: &str) {
    let listener = TcpListener::bind(address)
        .expect("failed to attach to {address}");

    let mut client_manager = ClientManager::new(listener);
    
    loop {
        client_manager.check_new_connections();
        client_manager.process_new_commands();
        client_manager.broadcast_to_subscribers( &heartbeat::get_heartbeat() );
        std::thread::sleep(HEARTBEAT_INTERVAL);
    }
}