import 'package:flutter/material.dart';
import 'package:auto_size_text/auto_size_text.dart';

import "custom_widget.dart";
import "globalData.dart";

class SelectServicePage extends StatelessWidget {
  const SelectServicePage({
    super.key,
    this.isAnAction = true,
  });
  final bool isAnAction;

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      theme: ThemeData(scaffoldBackgroundColor: const Color.fromARGB(255, 255, 255, 255)),
      home: APICall(isAnAction: isAnAction),
    );
  }
}

class APICall extends StatefulWidget {
  const APICall({
    super.key,
    this.isAnAction = true,
  });
  final bool isAnAction;

  @override
  SelectService createState() => SelectService(isAnAction: isAnAction);
}

class SelectService extends State<APICall> {
  SelectService({
    required this.isAnAction,
  });
  final bool isAnAction;
  final nbRepeat = 10;
  String title = "Select action service";

  @override
  Widget build(BuildContext context) {
    if (isAnAction != true) {
      title = "Select reaction service";
    }
    List<Widget> services = [];
    for (int i = 0; i < 5; i++) {
      services.add(
        Row(
          children: [
            SquareServicesCard(onTap: () {Navigator.pop(Globaldata.myContext);},),
            SquareServicesCard(onTap: () {Navigator.pop(Globaldata.myContext);},),
          ],
        ),
      );
    }
    return Scaffold(
      body: SingleChildScrollView(
        child: Column(
          children: [
            Padding(padding: EdgeInsets.only(top: 30)),
            SelectServiceTopBar(title: title),

            MyDivider(),
            Column(
              children: services,
            )
          ],
        ),
      ),
    );
  }
}

class SelectServiceTopBar extends StatelessWidget {
  const SelectServiceTopBar({
    super.key,
    required this.title,
  });

  final String title;

  @override
  Widget build(BuildContext context) {
    return Stack(
      alignment: Alignment.center,
      children: [
        Positioned(
          left: 0,
          child: IconButton(
            onPressed: () {Navigator.pop(Globaldata.myContext);},
            icon: Icon(Icons.arrow_back, size: 32),
          ),
        ),
        Center(
          child: AutoSizeText(title, style: TextStyle( fontSize: 30, fontWeight: FontWeight.bold, color: Colors.black,), maxLines: 1,),
        ),
      ],
    );
  }
}
