<!--

eyeOS Spice Web Client
Copyright (c) 2015 eyeOS S.L.

Contact Jose Carlos Norte (jose@eyeos.com) for more information about this software.

This program is free software; you can redistribute it and/or modify it under
the terms of the GNU Affero General Public License version 3 as published by the
Free Software Foundation.
 
This program is distributed in the hope that it will be useful, but WITHOUT
ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS
FOR A PARTICULAR PURPOSE.  See the GNU Affero General Public License for more
details.
 
You should have received a copy of the GNU Affero General Public License
version 3 along with this program in the file "LICENSE".  If not, see 
<http://www.gnu.org/licenses/agpl-3.0.txt>.
 
See www.eyeos.org for more details. All requests should be sent to licensing@eyeos.org
 
The interactive user interfaces in modified source and object code versions
of this program must display Appropriate Legal Notices, as required under
Section 5 of the GNU Affero General Public License version 3.
 
In accordance with Section 7(b) of the GNU Affero General Public License version 3,
these Appropriate Legal Notices must retain the display of the "Powered by
eyeos" logo and retain the original copyright notice. If the display of the 
logo is not reasonably feasible for technical reasons, the Appropriate Legal Notices
must display the words "Powered by eyeos" and retain the original copyright notice. 

-->

<script>
  const scripts = [
    '/lib/jquery-2.0.3.js',
    '/lib/jquery-mousewheel.js',
    '/lib/jgestures.min.js',
    '/lib/pixastic.js',
    '/lib/base64.js',
    '/lib/biginteger.js',
    '/lib/virtualjoystick.js',
    '/lib/prettyprint.js',
    '/lib/jsbn.js',
    '/lib/jsbn2.js',
    '/lib/prng4.js',
    '/lib/rng.js',
    '/lib/sha1.js',
    '/lib/encrypt.js',
    '/swcanvas/swcanvas.js',
    '/lib/bowser.js',
    '/lib/utils.js',
    '/lib/flipper.js',
    '/lib/CollisionDetector.js',
    '/lib/GlobalPool.js',
    '/lib/GenericObjectPool.js',
    '/lib/AsyncConsumer.js',
    '/lib/AsyncWorker.js',
    '/lib/PacketWorkerIdentifier.js',
    '/spiceobjects/spiceobjects.js',
    '/spiceobjects/generated/protocol.js',
    '/lib/graphicdebug.js',
    '/lib/images/lz.js',
    '/lib/images/bitmap.js',
    '/lib/images/png.js',
    '/lib/runqueue.js',
    '/lib/graphic.js',
    '/lib/queue.js',
    '/lib/ImageUncompressor.js',
    '/lib/SyncAsyncHandler.js',
    '/lib/IntegrationBenchmark.js',
    '/lib/stuckkeyshandler.js',
    '/lib/timelapsedetector.js',
    '/lib/displayRouter.js',
    '/lib/rasterEngine.js',
    '/lib/DataLogger.js',
    '/network/socket.js',
    '/network/clusternodechooser.js',
    '/network/socketqueue.js',
    '/network/packetcontroller.js',
    '/network/packetextractor.js',
    '/network/packetreassembler.js',
    '/network/reassemblerfactory.js',
    '/network/sizedefiner.js',
    '/network/packetlinkfactory.js',
    '/network/spicechannel.js',
    '/network/busconnection.js',
    '/network/connectioncontrol.js',
    '/application/agent.js',
    '/application/spiceconnection.js',
    '/application/clientgui.js',
    '/application/packetprocess.js',
    '/application/packetfilter.js',
    '/application/packetfactory.js',
    '/application/application.js',
    '/application/virtualmouse.js',
    '/application/imagecache.js',
    '/application/rasteroperation.js',
    '/application/stream.js',
    '/application/inputmanager.js',
    '/process/busprocess.js',
    '/process/displayprocess.js',
    '/process/displaypreprocess.js',
    '/process/inputprocess.js',
    '/process/cursorprocess.js',
    '/process/playbackprocess.js',
    '/process/mainprocess.js',
    '/keymaps/keymapes.js',
    '/keymaps/keymapit.js',
    '/keymaps/keymapus.js',
    '/keymaps/keymap.js',
    '/application/WorkerProcess.js',
    '/run.js'
  ];

  /*
  $effect.root(() => {
    for (let i = 0; i < scripts.length; i++) {
      const scriptElement = document.createElement('script');
      scriptElement.src = scripts[i];
      document.head.appendChild(scriptElement);
    }
  })
    */
</script>

<div id="testVdi"></div>
<input type="button" value="getStats" style="display:none;position:absolute;top:0px;right:0px;z-index:100;width:200px;height:200px;" id="getStats" />
<input type="text" id="hiddeninput" style="display:none;opacity:0;font-size: 60px;color:transparent"/>
<div id="canvasSpace" style="display: none; height: 1024px"></div>