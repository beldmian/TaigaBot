import { Navbar } from './components/Navbar.jsx'
import Container from '@material-ui/core/Container'

function App() {
  return (
    <div className="App">
      <Navbar/>
      <Container fixed>
        <h1>Taiga bot</h1>
      </Container>
    </div>
  );
}

export default App;
