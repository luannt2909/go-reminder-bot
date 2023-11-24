import React from 'react'
import {
    List,
    Datagrid,
    TextField,
    DateField,
    EditButton,
    DeleteButton,
    BooleanField
} from 'react-admin'

const ReminderList = (props) => {
    return (
        <List {...props}>
            <Datagrid>
                <TextField source='id' />
                <TextField source='name' />
                <BooleanField source='is_active' />
                <TextField source='schedule' />
                <TextField source='next_time' />
                <TextField source='webhook_type' />
                <TextField source='webhook' />
                <EditButton />
                <DeleteButton/>
            </Datagrid>
        </List>
    )
}

export default ReminderList