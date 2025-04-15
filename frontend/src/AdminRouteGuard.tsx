import *as React from 'react';
import {JSX} from 'react';
import {useAdminCheck} from "@/useAdminCheck";

export const AdminRouteGuard = ({children}: { children: JSX.Element }) => {
    const isAdmin = useAdminCheck();

    if (isAdmin === null) {
        // Показываем индикатор загрузки, пока идёт проверка
        return <div>Loading...</div>;
    }

    return isAdmin ? children : null;
};