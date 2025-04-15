import * as React from "react";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle
} from "@/components/ui/lib/alert-dialog.js";
import { useNavigate } from "react-router";

const CreateAccountDialogContext = React.createContext(null);

export const CreateAccountDialogProvider = ({ children }) => {
  const [open, setOpen] = React.useState(false);

  const navigate = useNavigate();

  const showDialog = () => setOpen(true);
  const hideDialog = () => setOpen(false);

  return (
    <CreateAccountDialogContext.Provider value={showDialog}>
      {children}
      <AlertDialog open={open} onOpenChange={setOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>An account needed!</AlertDialogTitle>
            <AlertDialogDescription>
              To continue, please, create free account
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel onClick={hideDialog}>Cancel</AlertDialogCancel>
            <AlertDialogAction onClick={() => {
              hideDialog();
              navigate("/register");
            }}>Create</AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </CreateAccountDialogContext.Provider>
  );
};

export const useCreateAccountDialog = () => {
  return React.useContext(CreateAccountDialogContext);
};
