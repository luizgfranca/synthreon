mod command;

mod client_manager;
mod application;
mod client;
mod heartbeat;


fn main() {
    application::serve_application("0.0.0.0:25250");
}
