import React from 'react'
import {BooleanInput, Create, SelectInput, SimpleForm, TextInput} from 'react-admin'
import {WebhookTypes} from './webhookType'

const TaskCreate = (props) => {
    return (
        <Create title='Create a Task' {...props}>
            <SimpleForm>
                <TextInput source='name' required fullWidth/>
                <BooleanInput source='is_active'/>
                <TextInput source='schedule' required fullWidth/>
                <SelectInput choices={WebhookTypes} required source='webhook_type'/>
                <TextInput multiline fullWidth required source='webhook'/>
                <TextInput multiline fullWidth required source='message'/>
            </SimpleForm>
        </Create>
    )
}

export default TaskCreate