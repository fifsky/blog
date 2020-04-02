<template>
    <div class="sect">
        <div id="blog-calendar">
            <table class="canlender">
                <tbody>
                <tr>
                    <td><a href="javascript:void(0);" @click="prevYear">&lt;&lt;</a></td>
                    <td><a href="javascript:void(0);" @click="prevMonth">&lt;</a></td>
                    <td colspan="3">{{currMonth}}</td>
                    <td><a href="javascript:void(0);" @click="nextMonth">&gt;</a></td>
                    <td><a href="javascript:void(0);" @click="nextYear">&gt;&gt;</a></td>
                </tr>
                <tr>
                    <td class="happy">日</td>
                    <td>一</td>
                    <td>二</td>
                    <td>三</td>
                    <td>四</td>
                    <td>五</td>
                    <td class="happy">六</td>
                </tr>
                <tr v-for="(days,index) in items" :key="index">
                    <td v-for="(d,i) in days" :key="i" :class="{today:isCurrDay(d)}">{{d}}</td>
                </tr>
                </tbody>
            </table>
        </div>
    </div>
</template>

<script>
    import dayjs from "dayjs"

    export default {
        name: "Calendar",
        data() {
            return {
                currDay: dayjs().format("YYYY-MM-DD"),
                currMonth: dayjs().format("YYYY-MM"),
                items: []
            }
        },
        methods: {
            prevMonth() {
                this.currMonth = dayjs(this.currMonth).subtract(1, 'month').format("YYYY-MM")
            },
            nextMonth() {
                this.currMonth = dayjs(this.currMonth).add(1, 'month').format("YYYY-MM")
            },
            prevYear() {
                this.currMonth = dayjs(this.currMonth).subtract(1, 'year').format("YYYY-MM")
            },
            nextYear() {
                this.currMonth = dayjs(this.currMonth).add(1, 'year').format("YYYY-MM")
            },
            builder() {
                let days = dayjs().daysInMonth()
                let startWeek = dayjs(this.currMonth + "-01").day()
                let dayItem = []
                for (let i = 0; i < startWeek; i++) {
                    dayItem.push('')
                }
                for (let i = 1; i <= days; i++) {
                    dayItem.push(i)
                }
                if (dayItem.length < 35) {
                    for (let i = 0; i < 35 - dayItem.length; i++) {
                        dayItem.push('')
                    }
                }
                let result = [];
                for (let i = 0; i < dayItem.length; i += 7) {
                    result.push(dayItem.slice(i, i + 7));
                }
                this.items = result
            },
            isCurrDay(d){
                return this.currMonth + '-' + d === this.currDay
            }
        },
        watch: {
            currMonth() {
                this.builder()
            }
        },
        created() {
            this.builder()
        }
    }
</script>