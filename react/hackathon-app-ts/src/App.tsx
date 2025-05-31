import React, { useState } from 'react';
import './App.css';
import { CustomHeader, headerButtons, SideBarButton } from './components/layout';
import { fireAuth } from './firebase';
import { onAuthStateChanged } from 'firebase/auth';
import { LoginForm } from './LoginForm';
import { LoginLayout } from './components/loginlayout';


function App() {
  const [loginUser, setLoginUser] = useState(fireAuth.currentUser);

  // ログイン状態を監視して、stateをリアルタイムで更新する
  onAuthStateChanged(fireAuth, user => {
    setLoginUser(user);
  });


  return (
    <>
      <LoginLayout>
        <LoginForm />
        {loginUser && <div>ログインしています</div>}
      </LoginLayout>
      <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
        <SideBarButton />
        <CustomHeader buttons={headerButtons} />
      </div>
    </>
  );
}

export default App;
