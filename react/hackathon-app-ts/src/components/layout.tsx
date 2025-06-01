import React from "react";
import styled from "styled-components";
import { Link } from 'react-router-dom';

type HeaderButtonType = {
    icon?: React.ReactNode;
    label: string;
    onClick?: () => void;
}
type SideBarButtonType = {
    label: string;
    onClick?: () => void;
    path: string;
}

type SideBarProps = {
    top: string;
    buttons: SideBarButtonType[];
}
type HeaderProps = {
    left?: string;
}
type CustomHeaderProps = {
    buttons: HeaderButtonType[];
}



export const SideBarButton: React.FC<{ buttons: SideBarButtonType[] }> = ({ buttons }) => {
    return (
        <SideBarContainer>
            {buttons.map((button, index) => (
                <StyledLink to={button.path} key={index}>
                    <SideBar top={`${index * 120}px`} buttons={[]}>
                        {button.label}
                    </SideBar>
                </StyledLink>
            ))}
        </SideBarContainer>
    );
};

export const CustomHeader: React.FC<CustomHeaderProps> = ({ buttons }) => {
    return (
        <HeaderContainer>
            {buttons.map((button, index) => (
                <HeaderButton
                    key={index}
                    onClick={button.onClick}
                    left={`${index * 170}px`}
                >
                    {button.icon && <span className="icon">{button.icon}</span>}
                    {button.label}
                </HeaderButton>
            ))}
        </HeaderContainer>
    );
};
//ヘッダー
const HeaderButton = styled.button <HeaderProps>`

    display: flex;
    align-items: center; 
    justify-content: center;
    left: ${props => props.left || 'auto'};
    height: 70px;
    width: 150px;
    background-color: #f0f0f0;
    padding: 10px 20px;
    border: none;
    border-radius: 20px;
    &:hover {
        background-color: rgb(24, 185, 226);
        color: #fff;
        cursor: pointer;
        transform: scale(1.05);
        transition: all 0.2s ease;
    }
    top: 20px;  
    z-index: 1000;
    overflow: visible;
`;

const HeaderContainer = styled.header`
    position: fixed;
    width: 100%;  
    height: 110px;  
    display: flex;
    flex-direction: row;  
    justify-content: left;  
    align-items: center;  
    padding: 20px;
    gap: 15px;
    left: 270px;
    top: 0;  
    background-color: #ffffff;  // 背景色を設定
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);  // 影を追加
`;

export const headerButtons: HeaderButtonType[] = [
    { label: "おすすめ", onClick: () => console.log("Header Button 1") },
    { label: "検索", onClick: () => console.log("Header Button 2") },
    { label: "通知", onClick: () => console.log("Header Button 3") },
    { label: "メッセージ", onClick: () => console.log("Header Button 4") },
    { label: "設定", onClick: () => console.log("Header Button 5") },
    { label: "プロフィール", onClick: () => console.log("Header Button 6") },
];


//サイドバー
const SideBar = styled.button <SideBarProps>`
    height: 75px;
    width: 180px;
    background-color: #f0f0f0;
    padding: 10px;
    border: none;
    border-radius: 20px;
    &:hover {
        background-color:rgb(24, 185, 226);
        color: #fff;
        cursor: pointer;
        transform: scale(1.05);
        transition: all 0.2s ease;
    }
    top: ${props => props.top || '0'};
    left: 100px;
    bottom: 0;
    z-index: 1000;
    overflow: visible;
`;
const SideBarContainer = styled.aside`
    position: fixed;
    height: 100vh; 
    width: 270px;
    display: flex;
    flex-direction: column;
    justify-content: flex-start;
    align-items: center;
    gap: 25px;
    margin-top: 150px;
    padding: 0;
    left: 0;
    top: 50%;
    transform: translateY(-50%);
    background-color: #ffffff;  
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1); 
    z-index: 1000;
`;

export const sideBarButtonPath: SideBarButtonType[] = [
    { label: "ホーム", path: "/" },
    { label: "検索", path: "/search" },
    { label: "通知", path: "/notifications" },
    { label: "メッセージ", path: "/messages" },
    { label: "設定", path: "/settings" },
    { label: "プロフィール", path: "/profile" },
];



const StyledLink = styled(Link)`
    text-decoration: none; /* リンクの下線を消す */
    color: inherit; /* 親要素の色を継承する */
`;
