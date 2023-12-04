import React from 'react'
import {
    BooleanField,
    CreateButton,
    DatagridConfigurable,
    DeleteButton,
    EditButton,
    List,
    TextField,
    TopToolbar,
} from 'react-admin'
import TestWebhookButton from "./TestWebhookButton";
import ContributeButton from "./ContributeButton";

// Usage
const ListActions = ({props}) => (
    <TopToolbar>
        <TestWebhookButton label="Webhook Test" {...props} />
        <CreateButton/>
        <ContributeButton/>
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