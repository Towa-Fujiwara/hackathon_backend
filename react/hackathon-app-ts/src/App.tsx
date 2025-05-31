import React, { useEffect, useState } from 'react';
import './App.css';
import { CustomHeader, headerButtons, SideBarButton } from './components/layout';
import { fireAuth } from './firebase';
import { onAuthStateChanged } from 'firebase/auth';
import { LoginForm } from './LoginForm';
import { LoginLayout } from './components/loginlayout';

const App = () => {
  const [loginUser, setLoginUser] = useState(fireAuth.currentUser);

  useEffect(() => {
    const unsubscribe = onAuthStateChanged(fireAuth, user => {
      setLoginUser(user);
    });
    return () => unsubscribe();
  }, []);


  return (
    <>
      {loginUser ? (
        // --- ログインしている場合に表示する内容 ---
        <>
          <SideBarButton />
          <CustomHeader buttons={headerButtons} />
          {/* 今後、タイムラインなどのメインコンテンツを
            このあたりに追加していくことになります。
          */}
        </>
      ) : (
        // --- ログインしていない場合に表示する内容 ---
        <LoginLayout>
          <LoginForm />
        </LoginLayout>
      )}
    </>
  );
}

export default App;
