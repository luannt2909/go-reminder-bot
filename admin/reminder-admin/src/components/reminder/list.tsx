import React from 'react'
import {
    BooleanField,
    Button,
    CreateButton,
    DatagridConfigurable,
    DeleteButton,
    EditButton,
    List,
    TextField,
    TopToolbar,
} from 'react-admin'
import {Link} from 'react-router-dom';
import IconGithub from '@material-ui/icons/Github';
import TestWebhookButton from "./TestWebhookButton";

const SourceCodeURL = "https://github.com/luannt2909/go-reminder-bot"

function ContributeButton({label}) {
    return <Link to={SourceCodeURL}>
        <Button label={label}>
            <IconGithub/>
        </Button>
    </Link>;
}

// Usage
const ListActions = ({props}) => (
    <TopToolbar>
        <TestWebhookButton label="Webhook Test" {...props} />
        <CreateButton/>
        <ContributeButton label="Contribute"/>
    </TopToolbar>
);
const ReminderList = (props) => {
    return (
        <List {...props}
              actions={<ListActions {...props}/>}
        >
            <DatagridConfigurable>
                <TextField source='id'/>
                <TextField source='name'/>
                <BooleanField source='is_active' label='Active'/>
                <TextField source='schedule'/>
                <TextField source='next_time'/>
                <TextField source='webhook_type'/>
                <TextField source='webhook' label='Webhook URL'/>
                <>
                    <EditButton/>
                    <TestWebhookButton label="Test" {...props}/>
                    <DeleteButton/>
                </>
            </DatagridConfigurable>
        </List>
    )
}

export default ReminderList