import { IconType } from "react-icons";
import styled from "styled-components";



type LoginLayoutProps = {
    children: React.ReactNode;
}
type LoginFormButtonProps = {
    onClick: () => void;
    label: string;
    icon: IconType;

}

export const LoginLayout: React.FC<LoginLayoutProps> = ({ children }) => {
    return (
        <LoginContainer>
            {children}
        </LoginContainer>
    );
};

const LoginContainer = styled.div`
    position: fixed;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    width: 500px; 
    height: auto; 
    padding: 40px;
    border-radius: 10px; 
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 20px; 
    background-color: rgb(24, 185, 226); 
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
`;

export const LoginButton: React.FC<LoginFormButtonProps> = ({ onClick, label, icon }) => {
    return (
        <StyledButton onClick={() => {
            onClick();
        }}>
            <Icon icon={icon} />
            <span>{label}</span>
        </StyledButton>
    );
};
const StyledButton = styled.button`
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: center;
    height: 50px;
    width: 150px;
    background-color: #ffffff;
    border-radius: 10px;
    border: none;
    cursor: pointer;
    &:hover {
        background-color: rgb(200, 196, 196);
        color: #fff;
        cursor: pointer;
        transform: scale(1.05);
        transition: all 0.2s ease;
    }
`;
const Icon = styled.div<{ icon: IconType }>`
    display: flex;
    align-items: center;
    justify-content: center;
    width: 20px;
    height: 20px;
`;

