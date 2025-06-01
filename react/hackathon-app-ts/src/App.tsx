import React, { useEffect, useState } from 'react';
import './App.css';
import { CustomHeader, headerButtons, SideBarButton, sideBarButtonPath } from './components/layout';
import { fireAuth } from './firebase';
import { onAuthStateChanged } from 'firebase/auth';
import { LoginForm } from './LoginForm';
import { LoginLayout } from './components/loginlayout';
import { BrowserRouter, Routes, Route, Outlet } from 'react-router-dom';
import { HomePage } from './pages/HomePage';
import { SearchPage } from './pages/SearchPage';
import { NotificationPage } from './pages/NotificationPage';
import { MessagePage } from './pages/MessagePage';
import { ProfilePage } from './pages/Profilepage';
import { SettingsPage } from './pages/SettingsPage';




const MainLayout = () => (
  <div>
    <SideBarButton buttons={sideBarButtonPath} />
    <CustomHeader buttons={headerButtons} />
    <main style={{ marginLeft: '270px', marginTop: '110px', padding: '20px' }}>
      {/* Outletは、子ルートのコンポーネントを描画するための場所 */}
      <Outlet />
    </main>
  </div>
);

const App = () => {
  const [loginUser, setLoginUser] = useState(fireAuth.currentUser);

  useEffect(() => {
    const unsubscribe = onAuthStateChanged(fireAuth, user => {
      setLoginUser(user);
    });
    return () => unsubscribe();
  }, []);


  return (
    <BrowserRouter>
      {loginUser ? (
        // --- ログインしている場合に表示する内容 ---
        <Routes>
          <Route path="/" element={<MainLayout />}>
            <Route index element={<HomePage />} />
            <Route path="/search" element={<SearchPage />} />
            <Route path="/notifications" element={<NotificationPage />} />
            <Route path="/messages" element={<MessagePage />} />
            <Route path="/settings" element={<SettingsPage />} />
            <Route path="/profile" element={<ProfilePage />} />
          </Route>
        </Routes>
      ) : (
        // --- ログインしていない場合に表示する内容 ---
        <LoginLayout>
          <LoginForm />
        </LoginLayout>
      )}
    </BrowserRouter>
  );
}

export default App;
