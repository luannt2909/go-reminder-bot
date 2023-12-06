import React from 'react'
import {Button, Link} from "react-admin";
import IconGithub from '@material-ui/icons/Github';

const SourceCodeURL = "https://github.com/luannt2909/go-reminder-bot"

const ContributeButton = ({...props}) => {
    return (
        <>
            <Link to={SourceCodeURL}>
                <Button label="Contribute">
                    {/*<IconGithub/>*/}
                </Button>
            </Link>
        </>
    )
}

export default ContributeButton;