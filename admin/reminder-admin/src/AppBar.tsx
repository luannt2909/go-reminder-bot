import {AppBar, IconButton, Toolbar, Button, Tooltip} from '@mui/material';
import {RefreshIconButton, SidebarToggleButton, TitlePortal, ToggleThemeButton, UserMenu} from 'react-admin';
import GithubIcon from '@mui/icons-material/GitHub';

const ContributeButton = () => (
    <Tooltip title="Github">
        <IconButton color="inherit" href="https://github.com/luannt2909/go-reminder-bot">
            <GithubIcon/>
        </IconButton>
    </Tooltip>
);

const LogoButton = () => (
    <Tooltip title="Luciango">
        <Button color="inherit" size='large'
                sx={{fontFamily: 'Cursive', fontSize: 'large', fontWeight: 'bold'}}>
            ~ Luciango ~
        </Button>
    </Tooltip>
);

export const CustomAppBar = () => (
    <AppBar color="primary">
        <Toolbar variant="dense">
            <>
                <SidebarToggleButton/>
                <TitlePortal/>
                <LogoButton/>
                <ContributeButton/>
                <ToggleThemeButton/>
                <RefreshIconButton/>
                <UserMenu/>
            </>
        </Toolbar>
    </AppBar>
);