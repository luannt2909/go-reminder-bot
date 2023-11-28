import React, {cloneElement} from 'react'
import {
    List,
    Datagrid,
    TextField,
    DateField,
    EditButton,
    DeleteButton,
    BooleanField,
    TopToolbar,
    SelectColumnsButton,
    CreateButton,
    DatagridConfigurable
} from 'react-admin'
import TestWebhookButton from "./TestWebhookButton";
const ListActions = ({props}) => (
    <TopToolbar>
        <TestWebhookButton label="Webhook Test" {...props} />
        <CreateButton/>
    </TopToolbar>
);
const ReminderList = (props) => {
    return (
        <List {...props}
        actions={<ListActions {...props}/>}
        >
            <DatagridConfigurable>
                <TextField source='id' />
                <TextField source='name' />
                <BooleanField source='is_active' label='Active' />
                <TextField source='schedule' />
                <TextField source='next_time' />
                <TextField source='webhook_type' />
                <TextField source='webhook' label='Webhook URL'/>
                <EditButton />
                <DeleteButton/>
            </DatagridConfigurable>
        </List>
    )
}

export default ReminderList