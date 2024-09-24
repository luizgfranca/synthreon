mod command;

mod client_manager;
mod application;
mod client;
mod heartbeat;


fn main() {
    application::serve_application("127.0.0.1:25250");
}
