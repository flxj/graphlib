/*
	Copyright (C) 2023 flxj(https://github.com/flxj)

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

		http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

package draw

const graphHTML = `<!DOCTYPE html>
<div id="container"></div>
<svg width="960" height="600"></svg>
<script src="https://cdnjs.cloudflare.com/ajax/libs/d3/7.9.0/d3.min.js"></script>
<script type="module">
    const width = 960;
    const height = 600;
    const obj = JSON.parse('{{.}}')
	const data = obj.data

    const simulation = d3.forceSimulation(data.nodes)
        .force("link", d3.forceLink(data.links).id(d => d.id))
        .force("charge", d3.forceManyBody())
        .force("center", d3.forceCenter(width / 2, height / 2));

    const svg = d3.select('svg')
      .attr("width", width)
      .attr("height", height)
      .attr("viewBox", [0, 0, width, height])
      .attr("style", "max-width: 100%; height: auto;");

    const link = svg.append("g")
        .attr("class", "links")
        .attr("stroke-opacity", 0.6)
        .selectAll("line")
        .data(data.links)
        .enter()
        .append("line")
		.attr("stroke", d => d.color!=""?d.color:"black")
        .attr("stroke-width",1.5);
    
    // TODO: add weight
    if (data.show_weight){
        const linkText = svg.append("g")
            .selectAll("text")
            .data(data.links)
            .enter()
            .append("text")
            .text(d => d.weight)
    }
    
    const node = svg.append("g")
        .attr("class", "nodes")
        .selectAll("circle")
        .data(data.nodes)
        .enter()
        .append("circle")
        .attr("r", 5)
        .attr("fill", d => d.color);
    // Add a drag behavior.
    node.call(d3.drag()
        .on("start", dragstarted)
        .on("drag", dragged)
        .on("end", dragended));
    
    // TODO 
    const nodeText = svg.append("g")
      .selectAll("text")
      .data(data.nodes)
      .enter()
      .append("text")
      .text(d => d.name)

    simulation.on("tick", () => {
        link
            .attr("x1", d => d.source.x)
            .attr("y1", d => d.source.y)
            .attr("x2", d => d.target.x)
            .attr("y2", d => d.target.y);

        node
            .attr("cx", d => d.x)
            .attr("cy", d => d.y);
    });
    //
    function dragstarted(event) {
        if (!event.active) simulation.alphaTarget(0.3).restart();
        event.subject.fx = event.subject.x;
        event.subject.fy = event.subject.y;
    }

    // Update the subject (dragged node) position during drag.
    function dragged(event) {
        event.subject.fx = event.x;
        event.subject.fy = event.y;
    }

    // Restore the target alpha so the simulation cools after dragging ends.
    // Unfix the subject position now that it’s no longer being dragged.
    function dragended(event) {
        if (!event.active) simulation.alphaTarget(0);
        event.subject.fx = null;
        event.subject.fy = null;
    }
     
    // When this cell is re-run, stop the previous simulation. (This doesn’t
    // really matter since the target alpha is zero and the simulation will
    // stop naturally, but it’s a good practice.)
    //invalidation.then(() => simulation.stop());
</script>
`

var digraphHTML string

const dotTemplate = `
{{if .Digraph }}digraph {{else}} graph {{end}}{
    {{range .Attr}}
    {{.}};
    {{end}}

    {{range .Nodes}}
    {{.}};
    {{end}}

    {{range .Edges}}
    {{.}};
    {{end}}
}
`
