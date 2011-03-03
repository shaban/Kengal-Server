<html style="font-family: Arial;">
<h1>Kengal</h1>
<h2>Distributed Blog System.</h2>
<h3>Features</h3>
<ul>
<li>Based on Golang http Package for being responsive</li>
<li>Using a minimal blog oriented subset of golang's template system to quickly adapt free templates</li>
<li>Out of the Box Master Slave replication via native mysql tools</li>
<li>centralized server and sitewide administration - TODO></li>
</ul>
<h3>Applications</h3>
<p>A Blog System structures its content as Category/Article Structure, while typical Articles do have a title , text, teaser, keywords etc. 
make this tool interesting for someone who needs customized templating for small amounts of articles which rules out a wiki approach.</p>
<h3>Usage</h3>
<pre>
	-h="": set Host MySql Adress like so -h myserver.com</li>
	-p="password": set Mysql Password for selected User here
	-u="root": set Mysql User here, default is root
	-db="mysql": set Database that MySql is supposed to connect to here</pre>
</html>
