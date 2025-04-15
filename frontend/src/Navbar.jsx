import { NavLink, useNavigate } from "react-router";
import { useAuth, useLogout } from "@/Auth.js";
import {
  Drawer,
  DrawerClose,
  DrawerContent,
  DrawerFooter,
  DrawerHeader,
  DrawerTitle,
  DrawerTrigger,
} from "@/components/ui/lib/drawer.js";
import { Separator } from "@/components/ui/lib/separator.js";
import { useDevice } from "@/useDevice.js";
import { useCreateAccountDialog } from "@/CreateAccountDialog.jsx";

export function Navbar() {
  const mobile = useDevice();

  if (mobile) return <MobileNavbar />;
  return <DesktopNavbar />;
}

function MobileNavbar() {
  const [isAuth, _, creds] = useAuth();
  const logout = useLogout();

  const navigate = useNavigate();

  const authError = useCreateAccountDialog();

  function DrawerLink({ children, link, auth = false }) {
    if (auth && !isAuth) {
      return <div onClick={() => authError()}>{children}</div>;
    }

    return (
      <DrawerClose asChild>
        <NavLink to={link}>{children}</NavLink>
      </DrawerClose>
    );
  }

  return (
    <div
      className={
        "border-t-solid fixed right-0 bottom-0 left-0 flex h-fit w-full flex-row flex-nowrap justify-end border-t-[2px] border-t-[zinc-600] bg-[red] bg-[rgba(0,0,0,0.5)]"
      }
      style={{ backdropFilter: "blur(4px)" }}
    >
      <Drawer>
        <DrawerTrigger asChild>
          <i className={"bi-list text-primary text-[4rem]"} />
        </DrawerTrigger>
        <DrawerContent>
          <div className="mx-auto w-full max-w-sm">
            <DrawerHeader>
              <DrawerTitle>
                <h1 className={"text-4xl"}>MetaCinema</h1>
              </DrawerTitle>
            </DrawerHeader>
            <div className="flex flex-col flex-nowrap gap-[0.6rem] p-4 pb-0 text-[2rem]">
              {isAuth ? (
                <DrawerClose asChild>
                  <div
                    onClick={() => {
                      logout();
                      navigate("/auth");
                    }}
                  >
                    <i className={"bi-box-arrow-right mr-[1em]"} />
                    {creds.username}
                  </div>
                </DrawerClose>
              ) : (
                <DrawerLink link={"/login"}>
                  <i className={"bi-box-arrow-in-right mr-[1em]"} />
                  Login
                </DrawerLink>
              )}
              <Separator />
              <DrawerLink auth link={"/couches"}>
                <i className={"bi-people-fill mr-[1em]"} />
                Couches
              </DrawerLink>
              <Separator />
              <DrawerLink link={"/popular"}>
                <i className={"bi-fire mr-[1em]"} />
                Popular
              </DrawerLink>
            </div>
            <DrawerFooter>
              <div
                className={
                  "flex flex-row flex-nowrap justify-between text-[2.5rem]"
                }
              >
                <DrawerLink auth link={"/myLikes"}>
                  <i className={"bi-heart"} />
                </DrawerLink>
                <DrawerLink auth link={"/explore"}>
                  <i className={"bi-compass"} />
                </DrawerLink>
                <DrawerLink link={"/search"}>
                  <i className={"bi-search"} />
                </DrawerLink>
              </div>
            </DrawerFooter>
          </div>
        </DrawerContent>
      </Drawer>
    </div>
  );
}

function DesktopNavbar() {
  const [isAuth, _, creds] = useAuth();
  const logout = useLogout();
  const navigate = useNavigate();

  const authError = useCreateAccountDialog();

  function NavbarLink({ children, auth = false, link, icon }) {
    if (auth && !isAuth)
      return (
        <div
          onClick={() => authError()}
          className="button flex items-center overflow-hidden px-4 py-2 text-2xl text-ellipsis whitespace-nowrap"
        >
          <i className={`${icon} mr-2 text-4xl`} />
          <span className="hidden lg:inline">{children}</span>
        </div>
      );

    return (
      <NavLink
        to={link}
        className="button flex items-center overflow-hidden px-4 py-2 text-2xl text-ellipsis whitespace-nowrap"
      >
        <i className={`${icon} mr-2 text-4xl`} />
        <span className="hidden lg:inline">{children}</span>
      </NavLink>
    );
  }

  return (
    <nav className="fixed top-0 right-0 left-0 flex w-full items-center justify-between border-b-2 border-zinc-600 bg-[rgba(0,0,0,0.5)] p-4 backdrop-blur-md">
      <NavLink to={"/"} className="text-primary button text-4xl">
        MetaCinema
      </NavLink>
      <div className="flex space-x-4">
        <NavbarLink link="/explore" auth icon="bi-compass">
          Explore
        </NavbarLink>
        <NavbarLink link="/popular" icon="bi-fire">
          Popular
        </NavbarLink>
        <NavbarLink link="/search" icon="bi-search">
          Search
        </NavbarLink>
      </div>
      <div className="flex space-x-6">
        <NavbarLink link="/couches" auth icon="bi-people-fill">
          Couches
        </NavbarLink>
        <NavbarLink link="/myLikes" auth icon="bi-heart">
          Likes
        </NavbarLink>
        {isAuth ? (
          <div
            onClick={() => {
              logout();
              navigate("/auth");
            }}
            className="button flex items-center px-4 py-2 text-2xl"
          >
            <i className="bi-box-arrow-right mr-2 text-4xl" />
            <span>{creds.username}</span>
          </div>
        ) : (
          <NavbarLink link="/login" icon="bi-box-arrow-in-right">
            Login
          </NavbarLink>
        )}
      </div>
    </nav>
  );
}
