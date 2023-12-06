import React, {useEffect, useState} from 'react'
import {TextInput} from "react-admin";
import {Cron} from 'react-js-cron'
import 'react-js-cron/dist/styles.css'
import {useFormContext} from 'react-hook-form';

const CronScheduleInput = props => {
    const [schedule, setSchedule] = useState('30 5 * * 1,6')
    const {setValue} = useFormContext();
    useEffect(() => {
        setValue('schedule', schedule)
    });
    const onChange = (event) => {
        setSchedule(event.target.value)
    }
    return (
        <>
            <TextInput source='schedule'
                       sx={{width: "50%"}}
                       required
                       onChange={onChange}
                       helperText="ex: '* * * * *', '@every 5m',... "/>
            <Cron value={schedule} setValue={setSchedule}/>
        </>
    )
}
export default CronScheduleInput;