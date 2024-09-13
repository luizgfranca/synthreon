import "./index.css";

const projects = ['Project A', 'Project B', 'Project C']

function App() {
  return (
    <div className="bg-zinc-900 text-zinc-100 h-screen">
      <div className="container mx-auto px-4 py-8">
        <h1 className="text-3xl font-bold mb-6">Projects</h1>
        
        <div className="space-y-4">
          {
            projects.map(project => (
              <div className="bg-zinc-800 hover:bg-zinc-700 p-4 cursor-pointer">
              <h2 className="text-xl font-semibold mb-2">{project}</h2>
              <p className="text-zinc-400">
                Long description of the project
              </p>
              </div>
            ))
          }  

          <div className="bg-zinc-900 hover:bg-zinc-700 p-4 text-center cursor-pointer">
            <i className="fas fa-plus mr-2 text-white"></i>Create New Project
          </div>
        </div>
      </div>
    </div>
  );
}

export default App;
