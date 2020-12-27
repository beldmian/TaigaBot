import { Navbar } from './components/Navbar.jsx'
import Container from '@material-ui/core/Container'

function App() {
  return (
    <div className="App">
      <Navbar/>
      <Container fixed>
        <h1>TaigaBot</h1>
        <p>
          Если вы когда-нибудь хотели простого бота, который бы выполнял все необходимые задачи, то это ваш выбор. Тайга поможет вам:
        </p>
        <ul>
          <li>Дать пользователям цветовые роли</li>
          <li>Хранить ваши задачи</li>
          <li>Удалять сообщения</li>
          <li>Выдавать роли для ивентов</li>
          <li>И конечно исать аниме</li>
        </ul>
        <p>Полный список всех возможностей бота вы можете узнать добавив его на свой сервер и написав <code>!help</code></p>
      </Container>
    </div>
  );
}

export default App;
