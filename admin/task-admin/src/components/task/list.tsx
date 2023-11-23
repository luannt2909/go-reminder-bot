import React from 'react'
import {
    List,
    Datagrid,
    TextField,
    DateField,
    EditButton,
    DeleteButton,
} from 'react-admin'

const TaskList = (props) => {
    return (
        <List {...props}>
            <Datagrid>
                <TextField source='id' />
                <TextField source='name' />
                <TextField source='schedule' />
                <TextField source='webhook_type' />
                <TextField source='webhook' />
                <EditButton />
                <DeleteButton/>
            </Datagrid>
        </List>
    )
}

export default TaskList