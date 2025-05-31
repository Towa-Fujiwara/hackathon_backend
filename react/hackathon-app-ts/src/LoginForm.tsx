import { signInWithPopup, GoogleAuthProvider, signOut } from "firebase/auth";
import { fireAuth } from "./firebase";
import { LoginButton } from "./components/loginlayout";
import { FcGoogle } from "react-icons/fc";

export const LoginForm: React.FC = () => {
    /**
     * googleでログインする
     */
    const signInWithGoogle = (): void => {
        // Google認証プロバイダを利用する
        const provider = new GoogleAuthProvider();

        // ログイン用のポップアップを表示
        signInWithPopup(fireAuth, provider)
            .then(res => {
                const user = res.user;
                alert("ログインユーザー: " + user.displayName);
            })
            .catch(err => {
                const errorMessage = err.message;
                alert(errorMessage);
            });
    };

    /**
     * ログアウトする
     */
    const signOutWithGoogle = (): void => {
        signOut(fireAuth).then(() => {
            alert("ログアウトしました");
        }).catch(err => {
            alert(err);
        });
    };

    return (
        <div style={{ display: 'flex', flexDirection: 'column', gap: '15px', alignItems: 'center' }}>
            <LoginButton onClick={signInWithGoogle} label="Googleでログイン" icon={FcGoogle} />
            <LoginButton onClick={signOutWithGoogle} label="ログアウト" icon={FcGoogle} />
        </div>
    );
};