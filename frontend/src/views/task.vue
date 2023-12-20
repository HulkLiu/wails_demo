<template>
  <n-data-table
      :columns="columns"
      :data="data"
      :pagination="pagination"
      :bordered="false"
  />
</template>

<script>
import { h, defineComponent } from "vue";
import { NButton, useMessage } from "naive-ui";

const createColumns = ({
     play
   }) => {
  return [
    {
      title: "No",
      key: "no"
    },
    {
      title: "Title",
      key: "title"
    },
    {
      title: "Length",
      key: "length"
    },
    {
      title: "Action",
      key: "actions",
      render(row) {
        return h(
            NButton,
            {
              strong: true,
              tertiary: true,
              size: "small",
              onClick: () => play(row)
            },
            { default: () => "Play" }
        );
      }
    }
  ];
};

const data = [
  { no: 3, title: "Wonderwall", length: "4:18" },
  { no: 4, title: "Don't Look Back in Anger", length: "4:48" },
  { no: 12, title: "Champagne Supernova", length: "7:27" }
];

export default defineComponent({
  setup() {
    const message = useMessage();

    return {
      data,
      columns: createColumns({
        play(row) {
          message.info(`Play ${row.title}`);
        }
      }),
      pagination: false,


    };
  }
});
</script>