<script lang="ts">
	import loadPDFWithEffect from "$lib/pdfjs";
</script>

<canvas use:loadPDFWithEffect={"../../static/sample-local-pdf.pdf"}></canvas>
