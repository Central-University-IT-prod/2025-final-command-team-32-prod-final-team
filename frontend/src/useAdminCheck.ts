import {useEffect, useState} from 'react';
import {useNavigate} from 'react-router-dom';
import {useAuth} from "@/Auth";
import {checkAdmin} from "@/client";

export const useAdminCheck = () => {
    const [isAdmin, setIsAdmin] = useState<boolean | null>(null);
    const [isAuth, , , token] = useAuth(true);
    const navigate = useNavigate();

    useEffect(() => {
        const verifyAdmin = async () => {
            try {
                await checkAdmin(token!);
                setIsAdmin(true);
            } catch (error) {
                if (error.message.includes("403")) {
                    setIsAdmin(false);
                    navigate('/'); // Перенаправляем при отсутствии прав
                }
            }
        };

        if (isAuth && token) {
            verifyAdmin();
        }
    }, [isAuth, token, navigate]);

    return isAdmin;
};