document.addEventListener('alpine:init', () => {
    Alpine.data('datepicker1', () => ({
        init() {
            const datepicker1 = M.Datepicker.init(this.$refs.datepicker1, {
                // Datepicker 1 options
            });
        }
    }));

    Alpine.data('datepicker2', () => ({
        init() {
            const datepicker2 = M.Datepicker.init(this.$refs.datepicker2, {
                // Datepicker 2 options
            });
        }
    }));
});
