import 'package:flutter/material.dart';
import 'package:auto_size_text/auto_size_text.dart';

import "custom_widget.dart";
import "globalData.dart";

class SaveNewWorkflowsPage extends StatelessWidget {
  const SaveNewWorkflowsPage({
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
  SaveNewWorkflows createState() => SaveNewWorkflows(isAnAction: isAnAction);
}

class SaveNewWorkflows extends State<APICall> {
  SaveNewWorkflows({
    required this.isAnAction,
  });
  final bool isAnAction;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: SingleChildScrollView(
        child: Column(
          children: [
            Container(
              decoration: BoxDecoration(
                color: Color.fromARGB(255, 63, 68, 69),
                border: Border.all(
                  color: Color.fromARGB(255, 63, 68, 69),
                  width: 2.0,
                ),
                borderRadius: BorderRadius.circular(4.0),
              ),

              child: InkWell(
                onTap: () {
                  debugPrint('Card tapped.');
                },
                child: Column(
                  children: const [
                    Padding(padding: EdgeInsets.only(top: 30)),
                    SelectServiceTopBar(title: "Review and finish", color: Colors.white),
                    Padding(padding: EdgeInsets.only(top: 20)),
                    LeftWidget(child: AutoSizeText("Workflow title", style: TextStyle( fontSize: 20, fontWeight: FontWeight.bold, color: Colors.white,), maxLines: 1,),),
                    WidgetCard (
                      backgroundColor: Colors.white,
                      borderRadius: 10,
                      child: CustomClassicForm(
                        hint:'Email',
                        padding: EdgeInsets.symmetric(horizontal: 0, vertical: 0),
                        focusedBorderColor: Colors.white,
                        borderWidth: 0,
                      ),
                    ),
                    LeftWidget(child: AutoSizeText("By {email/pseudo}", style: TextStyle( fontSize: 20, fontWeight: FontWeight.bold, color: Colors.white,), maxLines: 1,),),
                    Padding(padding: EdgeInsets.only(top: 10)),
                  ],
                ),
              ),
            ),

            Padding(padding: EdgeInsets.only(top: 30)),
            CustomButton (
              backgroundColor:Color.fromARGB(255, 34, 34, 34),
              text: AutoSizeText("Finish", style: TextStyle(fontSize: 30, fontWeight: FontWeight.bold, color: Colors.white),),
              onPressed: () {},
              size: Size(double.infinity, 60),
              padding: EdgeInsets.symmetric(horizontal: 15.0),
            ),
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
    this.color = Colors.black,
  });

  final String title;
  final Color color;

  @override
  Widget build(BuildContext context) {
    return Stack(
      alignment: Alignment.center,
      children: [
        Positioned(
          left: 0,
          child: IconButton(
            onPressed: () {Navigator.pop(Globaldata.myContext);},
            icon: Icon(Icons.arrow_back, size: 32, color: color,),
          ),
        ),
        Center(
          child: AutoSizeText(title, style: TextStyle( fontSize: 30, fontWeight: FontWeight.bold, color: color,), maxLines: 1,),
        ),
      ],
    );
  }
}
